package dao

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {
	col     *mongo.Collection
	liveCol *mongo.Collection
	node    *snowflake.Node
	ArticleDAO
}

func NewMongoDBArticleDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}

func (m *MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	art.Id = m.node.Generate().Int64()
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	_, err := m.col.InsertOne(ctx, art)
	return art.Id, err
}

func (m *MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	filter := bson.D{bson.E{Key: "id", Value: art.Id},
		bson.E{Key: "author_id", Value: art.AuthorId}}
	sets := bson.D{bson.E{Key: "$set",
		// 这里你可以考虑直接使用整个 art，因为会忽略零值。
		// 参考 Sync 中的写法
		// 但是我一般都喜欢显式指定要被更新的字段，确保可读性和可维护性
		Value: bson.D{bson.E{Key: "title", Value: art.Title},
			bson.E{Key: "content", Value: art.Content},
			bson.E{Key: "status", Value: art.Status},
			bson.E{Key: "utime", Value: time.Now().UnixMilli()},
		}}}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.MatchedCount != 1 {
		// 比较可能就是有人更新别人的文章，比如说攻击者跟你过不去
		return errors.New("没有修改权限，更新失败")
	}
	return nil
}

func (m *MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.UpdateById(ctx, art)
	} else {
		id, err = m.Insert(ctx, art)
	}
	if err != nil {
		return id, err
	}
	art.Id = id

	// UPDATE livecol 更新或插入
	filter := bson.D{bson.E{Key: "id", Value: art.Id},
		bson.E{Key: "author_id", Value: art.AuthorId}}
	now := time.Now().UnixMilli()
	art.Utime = now
	_, err = m.liveCol.UpdateOne(ctx, filter,
		bson.D{bson.E{Key: "$set", Value: ArticlePublish(art)},
			bson.E{Key: "$setOnInsert",
				Value: bson.D{bson.E{Key: "ctime", Value: now}}}},
		options.Update().SetUpsert(true))
	return id, err
}

func (m *MongoDBDAO) SyncStatus(ctx context.Context, author, uid int64, status uint8) error {
	filter := bson.D{bson.E{Key: "id", Value: uid},
		bson.E{Key: "author_id", Value: author}}
	sets := bson.D{bson.E{Key: "$set",
		Value: bson.D{bson.E{Key: "status", Value: status}}}}
	res, err := m.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return errors.New("没有修改权限，更新失败")
	}
	_, err = m.liveCol.UpdateOne(ctx, filter, sets)
	return err
}

package main

import (
	"code.google.com/p/gorest"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

type Quote struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Body string
	Score int
	Date time.Time
}

var (
	mgoSession *mgo.Session
	databaseName = "qdb"
)

func getSession() *mgo.Session {
	if (mgoSession == nil) {
		var err error
		mgoSession, err = mgo.Dial("127.0.0.1")
		if err != nil {
			panic(err)
		}
	}
	return mgoSession.Clone();
}

func main() {
	gorest.RegisterService(new(QuoteService))
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":8080",nil)
}

type QuoteService struct {
	gorest.RestService	`root:"/quotes-service/" consumes:"application/json" produces:"application/json"`

	listQuotes gorest.EndPoint `method:"GET" path:"/quotes" output:"[]Quote"`
	addQuote gorest.EndPoint `method:"POST" path:"/quotes" postdata:"Quote"`
	upVote gorest.EndPoint `method:"OPTIONS" path:"/quotes/{ID:string}/upvote"`
	downVote gorest.EndPoint `method:"OPTIONS" path:"/quotes/{ID:string}/downvote"`
}

func (serv QuoteService) ListQuotes() []Quote {
	s := getSession()
	defer s.Close()

	quoteCollection := s.DB(databaseName).C("quotes")

	var qtLst []Quote
	err := quoteCollection.Find(bson.M{}).All(&qtLst)
	if err != nil {
		panic(err)
	}


	return qtLst;
}

func (serv QuoteService) AddQuote(q Quote) {
	s := getSession()
	defer s.Close()

	quoteCollection := s.DB(databaseName).C("quotes")

	q.Score = 0;
	q.Date = time.Now();

	err := quoteCollection.Insert(&q)
	if err != nil {
		panic(err)
	}

	serv.ResponseBuilder().Created("/quotes-service/quotes/"+string(q.ID)) 
}

func (serv QuoteService) UpVote(ID string) {
	s := getSession()
	defer s.Close()

	quoteCollection := s.DB(databaseName).C("quotes")

	var returnedQuote Quote
	change := mgo.Change{
		Update: bson.M{"$inc": bson.M{"score": 1}},
		ReturnNew: true,
	}
	_, _ = quoteCollection.FindId(bson.ObjectIdHex(ID)).Apply(change, &returnedQuote)
}

func (serv QuoteService) DownVote(ID string) {
	s := getSession()
	defer s.Close()

	quoteCollection := s.DB(databaseName).C("quotes")

	var returnedQuote Quote
	change := mgo.Change{
		Update: bson.M{"$inc": bson.M{"score": -1}},
		ReturnNew: true,
	}
	_, _ = quoteCollection.FindId(bson.ObjectIdHex(ID)).Apply(change, &returnedQuote)
}
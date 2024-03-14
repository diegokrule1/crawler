# crawler

This application crawls a web app base on an initial url

## Compilation

This is a golang application so you need `go` to compile it.

Then run `go build` to compile it

## Design decisions

There is a database to stored already visited urls. The database is postgres
and there is only one table which is defined on [here](https://github.com/diegokrule1/crawler/blob/start/migrations/url_up.sql)

The application is design following a producer-consumer patttern, where the producer stores
the newly found url on the database and the consumer does the analysis.
Since there is a `constraint unique` on the url on the database we make sure 
never to analyze the same url twice.

The fact that we are using channels and go routine allows the application to
analyse several url concurrently

## Limitations

The current crawler only "walks" the url. Much more could be done by analysing the responses.
Currently only "text/html" responses are considered. This was only to keep this simple.

Also on the database the structure winds up being a tree and not a more complex graph. That is,
a page only has one father. This does not always hold, you can get to one page from different places.
This is not being captured on the current design just for simplicity reasons.

In order for the app to come to an end there in an internal cronjob, in the form
of a ticker, that is periodically inspecting the database in order to see if there is any work 
left to be done.
Urls have state on the database. When one is created it is on the state "created", when it is being analysed it
is on the state "processing". Finally when the url analysis is finished it is on the state "processed".
The app goes on as long as there is at least one url on a state which is not "processed"
This can be further improved by having a notification mechanism from the producers to the consumers.

Only one integration test has been made to test the happy path case. This is
not enough to ensure the application correctness. More test need to be added 
which test more cases

## How To Run
Assuming you build the application using `go build`, a binary file called `crawler` will be generated.
To run the application run

`crawler crawler {url}`
 
where `url` is the url you want to crawl 
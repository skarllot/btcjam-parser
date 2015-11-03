# btcjam-parser
A BTCjam transaction parser for local querying

# Usage

Log in BTCjam and open the following link: [https://btcjam.com/transactions.json?iDisplayStart=0&iDisplayLength=1000000000&type=all][https://btcjam.com/transactions.json?iDisplayStart=0&iDisplayLength=1000000000&type=all] in your browser.

Save the content to file (eg: btjam-transactions.json).

Parse to returned JSON to a new sane format:

~~~ bash
./btcjam-parser < btcjam-transactions.json > sane-transactions.json
~~~

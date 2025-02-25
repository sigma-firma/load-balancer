# load-balancer (lb)

`lb` is an MIT licensed application written in `go` designed to simulate a load balancer

[Screen Recording 2025-02-24 at 7.59.22 PM.webm](https://github.com/user-attachments/assets/b0ebaf5f-bcf9-4105-90e8-4175ede5bd6f)

This program acts a model "load balancer", accepting incoming "requests"
at a rate moderated via the channel provided via the `go` `*time.Ticker` type
(`<-&t.T{}.C`) and responding to those requests via the same approach.
As requests come in, the number of open connections are balanced determined
accordingly via the request rate flag (`--req_rate`) and response time flag
(`--res_rate`), over a number of servers in different (metaphorical) 
regions (`--regions`), starting with the next "closest" server. Each servers 
maximum number of connections is set via the max conns flag (`--max_conns`), 
thus the max cluster connections can be determined via the number of regions
multiplied by the max conns. 

To run the program we use a command like this:

       $ ./lb --regions=14 --req_rate=13 --res_rate=200 --max_conns=9

         --regions  ->  the number of regions (servers) to simulate.
        --req_rate  ->  the number of request to send per millisecond.
        --res_rate  ->  the number of milliseconds it takes to handle each
                        request.
       --max_conns  ->  the maximum number of connection a server can have.

As these parameters are adjusted, one may notice that some configurations
overload the cluster near instantly, while some may do so more gradually,
and, if one is lucky (I don't need luck) one may find sustainable
configurations, capable of maintaining the clusters indefinite integrity.

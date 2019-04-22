## Run
run urlcrawler_test.go or run urlcrawler_test.go with coverage

## Approach to Design 

In the above question, we have only been asked to crawl 2 levels for simplicity & avoid huge amounts of data, however, if we designed a solution that could only crawl 2 levels, our solution would not be able to scale to any number of level/depth_n. 
Tailoring a solution to only depth_2, we could have taken a basic logic approach(only store and traverse internal urls in origin). Only output external urls in html bodies of (origin + internal urls in origin).

Typically, every domain has multiple internal pages and every internal page has sub internal pages. For example-- example.com/abc, example.com/about, example.com/static, example.com/abc/originstatic.

+ **Data structures**

A user may wish to crawl all internal pages of a domain or to any given depth_n. So we'd have to come with an approach that allows to keep track of levels & identify such dependencies in current internal page to previously visited internal pages or unvisited internal page while traversing each page-- How do we track all internal links/dependencies and store the state of the system? Thinking of this entire system as a graphs elegantly maps the ancestor-descendant dependencies with directed edges. In fact, a graph without any cycles is a tree, so we can think of this as a tree since we should not create any cycles.

Let’s think of every url-path as a node and every sub-internal-url-path as a child/successor node. There is a directed edge from parent-internal path to child sub-internal-path. If there are two distinct url-paths that refer to the same HTML content, then we should skip traversing the duplicate to avoid cycles in the graph. 

+ **Algorithm** 

How do we discover and add nodes/edges into our graph?

We can either do DFS or BFS approach to this - since the order in which we output our external URLs does not matter. 

Although, doing DFS may cause our memory stack to blow up when there are many levels. Also thinking of parallelization, DFS/any-top-down-traversal will cause a lot of idly waiting threads. Let’s say we do divide and conquer, run DFS and spawn new goroutines for nodes at level_i (** Limited by ideal maximum number of goroutines ideal for depth_n and that would not cause the hit on performance) Ancestors nodes depend on their child nodes to return values to proceed, in effect, many goroutines would be idly waiting. Ideally, we would not want goroutines to be idly waiting. 

This means it would be useful to process a node(each HTML body) fully and extract all internal nodes(internal-sub-paths) before we proceed to next thus removing ancestor-descendant dependency(repeating until depth_n). Also, this allows to run parallel computations with proper sync and concurrency on each node without waiting.

This hints at level-by-level or BFS traversal, we can push each unvisited node(internal-sub-paths) onto a queue to process it in the order it arrived(FIFO), we keep adding nodes into the queue at every level before de-queueing and processing the next descendant level. 

Since the problem doesn't ask for visualization or to output the entire state & dependencies of paths and subpaths we don’t need to maintain an entire graph(adjacency lists), instead we can just do a BFS on it with queues(+ other states)

Note- each of these may have any file extension(html, htm, php, asp) The default homepage/HTMLBody(and for that matter any internal page) can be referred to again with non-unique paths that haven't been seen leading to same HTMLBody, so we would need a way to avoid such cycles.


## Scalablibility and further thoughts

With the above approach we can traverse any depth_n instead of only until given depth_i. Since we can run parallel computations it allows us to distribute our load on multiple cores thus reducing latency. This also means we'd have to maintain shared memory which is concurrent allowing for concurrent reads/writes by multiple goroutines(threads). (Note- how using a mutex on a non-concurrent data structure can stall performance). 

This problem limits us to only internal URLs, however, we can easily think of a scenario where the user may want to crawl external URLs as well in the process(after all those are also domain names and have multiple internal urls). 

Also, at scale, we would typically run multiple instances of URL Crawler on multiple servers. Instead of keeping a store/amp per instance a global concurrent store/map (with replicas) would make more sense in order to avoid for duplicate crawling. 

We would have to make sure our system is failure tolerant, so if a server crashes we should not lose the state of the system. So, we should keep persistent store (distributed key-value store which has high performance)

These multiple instances can run in parallel and each instance would have parallelism within in allowing for low latency and faster traversals.

## Deployment Pipeline

As the system grows larger, we can think of breaking up the components into smaller microservices (for example- GetStatus, GetResult, CrawlDomainURL Monitoring/Logging). Each of these instead of running in VMs or bare-metal can run in docker containers because those are enviornment agnostic, light-weight, faster-deployment. 

Since we could have many instances of the same URLCrawler, we would have lots of docker containers and to manage netowrking/storage/communication-across-components/bootstrapping-the-system/rollbacks/upgrades/loadbalancing-user-traffic we'd rather use container orchestration system for automating application deployment, scaling, and management- so using Kubernetes allow us to manage our docker containers in pods and handle availbility/reliability of the system.

We can run this entire system in first party data centers or in public cloud which takes away the burden of managing our own hardware.

## Further Optimizations/Performance Evaluation

The HTMLBodies of each page can be huge- So we can also parallelize this by dividing/conquering the HTMLBodies into independent parts and extracting internal/external URLs.

Also GetResult can be optimized by using a map instead of an array, and if order does not matter, we can split computations of the result ds, and traverse the ds faster.

For small HTMLBodies and relatively small depths, spawning many goroutines may infact stall performance since goroutine scheduling and management(waits/locking/syncs/pools) is a (computationally+memory)-intesive task.

So we must spawn goroutines based on depth and size of the data to balance the trade off between computational latency and goroutine management latency.

## Current Limitations and Future Work

Although I caught errors at every point where errors were thrown, and return values as appropriate, however, upon debugging found that HTTP get on some URLs with complex response bodies times out due to various reasons and wait endlessly on certain URLs- found some clues online why this may happen- sending invalid authorization credentials, apis returning invalid content type when there is an error, using wrong http method, not handling unexpected error codes properly so this is a area where we can make this system way more robust to be able to work on any given URL and catch for all possible cases that may happen.

Being able to provide visualization of the entire graph with dependencies would be very interesting.

A lot of above optimizations/scaling/deployment-ways can be implemented to make the system reliable, available and consistent.

Parsing every href 'a' tag can be also done more robustly since currently we also encounter a lot of javascipts/mail-to/other-non-URL parts in it.

Malformed URLs should directly be handled from http.get requests and checking for HTTP Response header to ensure that the body is well-formatted, else return appropriately. We can extend the MimeTypeSet to include all possible types(list of known mimeTypes- https://www.freeformatter.com/mime-types-list.html). Even After this the response body can be corrupted, so those can be better handled with more logic.



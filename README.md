## Run
run urlcrawler_test.go or run urlcrawler_test.go with coverage

## Approach to Design 

In the above question, we have only been asked to crawl 2 levels for simplicity & avoid huge amounts of data, however, if we designed a solution that could only crawl 2 levels, our solution would not be able to scale to any number of level/depth_n. 

Imagining tailoring a solution to only depth_2, we could have easily taken a basic logic approach(only store and traverse internal urls in origin). Only output external urls in html bodies of (origin + internal urls in origin).

Typically, every domain has multiple internal pages and every internal page has sub internal pages. For example-- example.com/abc, example.com/about, example.com/static, example.com/abc/originstatic.

A user may wish to crawl all internal pages of a domain or to any given depth_n. So we would have to come with an approach that allows us to identify dependencies in current call to internal page to previously visited internal pages or unvisited internal page while going each page, so question becomes how do we track all internal links/dependencies and store the state of the system? Thinking of this entire system as a graphs elegantly maps the ancestor-descendant dependencies and directed edges that graphs offer. In fact, a graph without any cycles is a tree, so we can think of this as a tree since we should not create any cycles.

Let’s think of every url-path as a node and every sub-internal-url-path as a child/successor node. There is a directed edge from parent-internal path to child sub-internal-path. If there are two distinct url-paths that refer to the same HTML content, then we should skip traversing the duplicate to avoid cycles in the graph. 

Next, how do we discover and add each node into our graph, edges and build out our graph?

We can either of DFS or BFS approach to this - since the order in which we output our external URLs does not matter. 

Although, doing DFS may cause our memory stack to blow up when there are many levels. Also thinking of parallelization, DFS/any-top-down-traversal will cause a lot of idly waiting threads. Let’s say we spawn new goroutines for nodes all level_i (*Limited by ideal maximum number of goroutines ideal for depth_n and that would not cause the hit on performance) Ancestors nodes depend on their child nodes to return values to proceed, in effect, many goroutines would be idly waiting. Ideally, we would not want goroutines to be idly waiting. 

This means it would be useful to process node(each HTML body) fully and extract all internal nodes(internal-sub-paths) and then run parallel computations with proper sync and concurrency on each node independently and remove any ancestor-descendant dependency, repeating the same process over until we reach depth_n.

This hints at level-by-level or BFS traversal where-in we can push each unvisited node(internal-sub-paths) onto a queue to process it in the order it arrived(FIFO), we keep adding every nodes into the queue that we encounter at every level before moving onto de-queueing and processing the next descendant level. 

Since the problem does not ask us for visualization or to output the entire state & dependencies of paths and subpaths we don’t necessarily need to maintain an entire graph in some data structure(adjacency lists), instead we can just do a BFS on it with queues and additional structures.

Note- each of these may have any file extension(html, htm, php, asp) The default homepage(and for that matter any internal page) can be referred to again with non-unique paths already seen for the same page leading the same HTML body, so we would need a way to avoid such cycles.


## Why is this scalable?

With the above approach we can traverse any depth_n instead of only till a given depth_i. Since we can run parallel computations it allows us to distribute our load on multiple cores thus reducing latency. This also means we would have to maintain a shared memory that must be concurrent allowing for concurrent reads and writes by multiple goroutines(threads) simultaneously. (Note- how using a mutex on a non-concurrent data structure can stall performance). 

This problem limits us to only internal URLs, however, we can easily think of a scenario where the user may want to crawl external URLs as well in the process(after all those are also domain names and have multiple internal urls). 

Also, at scale, we would typically run multiple instances of URL Crawler on multiple servers. Instead of keeping a storage per instance a global concurrent storage (with replicas) would make more sense in order to avoid for duplicate crawling. 

We would have to make sure our system is failure tolerant, so if a server crashes we should not lose the state of the system. So, we should keep persistent storage 

These multiple instances can run in parallel and each instance would have parallelism within in allowing for low latency and faster traversals.

## Deployment Pipeline

## Further Optimizations/Performance Evaluation

## Current Limitations and Future Work

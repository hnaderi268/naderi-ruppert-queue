# naderi-ruppert-queue
``naderi-ruppert-queue`` is a shared data strucutre based on a paper from [Naderi & Ruppert](https://github.com/hnaderi268/naderi-ruppert-paper) (ACM PODC 2023). A shared data structure is a data structure where multiple processes can perform operations on it at the same time. naderi-ruppert-queue is linearizable, meaning that operations appear to take place atomically, and wait-free, meaning that some operation on the queue is guaranteed to complete regardless of how asynchronous processes are scheduled to take steps.

A binary tree, called the block-tree  is used to store the sequence of operations in ``nader-ruppert-queue``. Each process has a leaf in the block-tree where it adds its operations to it. Operations are propagated from the leaves up to the root in a cooperative way that ensures wait-freedom and avoids the CAS retry problem. Operations in the root are ordered, and this order is used to linearize the operations and compute their responses. In this repository, we have implemented a space-bound version of the algorithm.

# Code Strucutre

``algorithm`` package is divided into three files. ``block-tree.go`` is the primary data structure.
``queue.go`` includes the Enqueue and Dequeue functions. It's important to note that each process should perform only one operation at a time.
``help.go`` consists of various helper methods for performing Compare And Swap (CAS) operations.

In addition, the ``io`` package contains code for visualizing the block-tree and the order in which operations start and finish. This visualization is helpful for demonstrating linearization points.


# Usage

You can create a queue supprting 4 processes and 100 total operations like this.
```
q := algorithm.NewBlockTree(4, 100)
```
Then you can share ``q`` among processes ``P0, P1, P2, P3`` and have each process run ``Enqueue(q, v, pid)`` and ``Dequeue(q, pid)``. Here's an example of four processes performing these operations, where the time axis goes from top to bottom.
```
        P0        |        P1        |        P2        |        P3  
------------------+------------------+------------------+------------------
Enqueue(q,x,0)    |Enqueue(q,y,1)    |                  |Dequeue(q,3)
                  |Dequeue(q,1)      |                  |
Dequeue(q,0)      |                  |Enqueue(q,z,2)    |

```

# Run

There is a sample execution in ``main.go`` which you can run by ``go run .``. It also prints a visualization of the ordering of operations. ``[i`` shows the start of an operation by process ``i`` and ``i]`` shows the end of it. ``v -> ▢▢▢`` is for showing enqueuing ``v`` and ``▢▢▢ ->`` shows dequeuing ``v``. Example of a run on three processes 0,1,2:
````
0123 | PID |         OP
-----+-----+-------------------
  ▖  | [2  | 
▖ ▌  | [0  |  01 -> ▢▢▢
▌▖▌  | [1  |  11 -> ▢▢▢
▘▌▌  |  0] |    
▖▌▌  | [0  |  02 -> ▢▢▢
▘▌▌  |  0] |    
▖▌▌  | [0  | 
▌▌▘  |  2] |        ▢▢▢ -> 01
▌▌▖  | [2  |  21 -> ▢▢▢
▌▌▘  |  2] |    
▌▘   |  1] |    
▌▖   | [1  | 
▘▌   |  0] |        ▢▢▢ -> 02
 ▘   |  1] |        ▢▢▢ -> 11
 ▖   | [1  |  12 -> ▢▢▢
 ▘   |  1] |    
````

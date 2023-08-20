# naderi-ruppert-queue
Wait-free Queue with Poly-logarithmic Worst-case Step Complexity
based on paper from Naderi & Ruppert - ACM PODC 2023

The algorithm  stores the elements of the queue into a binary tree which we call block-tree. Each process has a leaf in the block-tree that stores the ordering of the operations(Enqueues and Dequeues) that it will call. Each node in the tree contains the ordering of the operations in its subtree. When process p wants to do operation op, it appends op in to its leaf's ordering and then propagates op up to the root. When op reaches to the root, p computes op's result (if op is an Enqueue). To solve the race conditions caused by multiple processes writing into one node of the block-tree we use CAS(Compare&Swap).

To test the algorithm you can run the project by ``go run .``. There is a sample execution in ``main.go`` which you can change.

``algorithm`` directory contains the main logic separated into three files. ``block-tree.go`` contains the main data-structure. ``queue.go`` contains ``Enqueue(v, pid)`` and ``Dequeue(pid)``. Note that each process should do only one operation at a time. ``help.go``  includes some helper methods for CAS and some small stuff. ``io`` package contains the code to visualize the tree and the ordering that operations started and finished which is useful for  showing linearization points.

``[i`` shows the start of an operation by process ``i`` and ``i]`` shows the end of it. ``v -> ▢▢▢`` is for showing enqueuing ``v`` and ``▢▢▢ ->`` shows dequeuing ``v``. Example of a run on three processes 0,1,2:
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
Note in te run above 02 is finished before 11, so it is dequeued earlier than that. Linearizability is more discussed in the paper.
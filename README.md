# B-Tree

## Overview

A simple implementation of a B-tree in Go, including:
- Insertion
- Deletion
- Searching

## What is a B-Tree?

A B-tree is a data structure commonly used in databases. It allows faster lookup for large datasets, especially compared to a binary tree, because the tree is shallower. However, for small datasets, it might be slower due to a higher number of comparisons.

In a B-tree:
- Every node can have up to M keys and M+1 children.
- Every node (except the root) must have at least ceil(M/2) - 1 keys.
- The root can have as few as zero or one key.
- Insertion of new keys must maintain the balanced structure of the tree.

This structure ensures efficient performance for operations like searching, inserting, and deleting.


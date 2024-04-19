## language options

I'm probably going to use go.
Normally, garbage collection would be a problem with a dbms (although i think cassandra or something is written in java? so it's not unheard of) but I think the simplicity will make it way easier to actually finish the project and that's worth the performance tradeoffs.

## file format & compatibility

I don't feel super strongly about maintaining compatibility with SQLite, but reusing their file format means I don't have to decide on a file format with which to store the necessary data for a database.

### hot journals

while the database state is usually contained in a single file on disk, additional information can be stored in a second file.
If the database host crashes during a transaction, this second file is called a "hot journal" and contains information vital to recovering the state of the database.

### pages

SQLite database files are made up of one or more _pages_[^1] .
The size of a page in bytes is 2^x where x satisfies the range [9, 16].
The maximum page number is 2^32 - 2 (4,294,967,294) meaning that most systems will run out of disk space long [^2] before a database reaches maximum capacity.

[^1]: pages and files are not interchangeable. a fie contains multiple pages. a database mostly consists of a single file.
[^2]: a SQLite databse can store up to about 281 terabytes of data

Each _page_ has **exactly one** use at each point in time:

- The lock-byte page
- A freelist page
  - A freelist trunk page
  - A freelist leaf page
- A b-tree page
  - A table b-tree interior page
  - A table b-tree leaf page
  - An index b-tree interior page
  - An index b-tree leaf page
- A payload overflow page
- A pointer map page

### database header

the database header format can be found [here](https://www.sqlite.org/fileformat2.html#the_database_header). No point copy-pasting it.

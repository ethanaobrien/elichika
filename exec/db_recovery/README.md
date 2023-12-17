# Recover db
Recover the deleted things from master database.

## How

There are database diff projects for [Global](https://github.com/kotori2/-llasww_master_db_diff) and [Japanese](https://github.com/kotori2/llas_master_db_diff) server by kotori2 (CPPO).

To get a specific patch, it's as easy as checking out the correct patch and build the database from scratch. However, we want to add deleted content back, so we will need to do some work.

### Generating the SQL
We use 
```
git log -p --unified=0
```
to generate the diff, then parse the rows that start with `"+INSERT INTO"` and `"-INSERT INTO"`.

We will keep the one that was competely removed, to be added onto the latest official version. However, we might also do some modification too. For example, sometime a thing is removed by marking it as removed and not be removing it from the database.

### Fixing the scripts
Some SQL commands will not be entirely correct. This is due to different in behaviour of SQL implementation (probably).

For example:
- Some NULL field are described as `""` instead of `NULL`.
- Some string fields have the wrong escape format, using `\"` instead of `""` for escaping `"`.

### Add things and maybe other things
When we have the correct script, we should run it, but we might run into errors.

Most of the time it will be in the form of a FOREIGN KEY violation. This is either caused by `NULL` field not being treated as `NULL`, or by the referenced stuff being missing. This has to be resolved case by case basis.

## Usage
This package implement the parsing part, by defining relevant rules in Golang, and output SQL that should be applied. You will have to check the implement itself.

For now, this package is only used to make elichika run to the intended degree.


```
根据实际的文件目录，采用深度优先搜索指定的目录．并在合适的位置
安置$guard-file文件．
注意:所有-开头的选项必须在指定目前之前设置.

  -clean-guard
    	clean the guard files and exit
  -count int
    	number of guard file you wish to set. The realy guard files is not exactly same as the number (default 100)
  -debug
    	
  -guard-file-name string
    	the guard file name (default "__GUARD__1456207797")
  -index-file-name string
    	the guard index name under root directory (default "__GUARD__INDEX__")
  -index-url string
    	the http url of index file, it should be used with report-mirror-progress
  -report-sync-progress string
    	the http url of mirror to be reported

```

Generate Guard files in server
```
export $repo = /data/www/mirrors/deepin
./repo_guard -clean-guard=true $repo/pool $repo/dists

./repo_guard $repo/pool $repo/dists
```

Report the sync progrss in anywhere
```
 % ./repo_guard -report-mirror-progress=$MirrorServer -index-url=$IndexURL
```

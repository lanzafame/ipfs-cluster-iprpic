#!/usr/bin/env fish

for i in (seq 1 50)
    echo $i | ipfs add -q | xargs ipfs-cluster-ctl pin add
end



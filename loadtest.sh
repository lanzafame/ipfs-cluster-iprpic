#!/usr/bin/env fish

for i in (seq 100 200)
    echo $i | ipfs add -q | xargs ipfs-cluster-ctl -l /ip4/192.168.110.63/tcp/9094 pin rm
end



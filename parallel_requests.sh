#! /usr/bin/bash
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request1" -X POST http://localhost:3000/shorten && echo "done1" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request2" -X POST http://localhost:3000/shorten && echo "done2" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request3" -X POST http://localhost:3000/shorten && echo "done3" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request4" -X POST http://localhost:3000/shorten && echo "done4" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request5" -X POST http://localhost:3000/shorten && echo "done5" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request6" -X POST http://localhost:3000/shorten && echo "done6" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request7" -X POST http://localhost:3000/shorten && echo "done7" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request8" -X POST http://localhost:3000/shorten && echo "done8" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request9" -X POST http://localhost:3000/shorten && echo "done9" &
curl -i -d "https://en.wikipedia.org/wiki/Main_Page/request10" -X POST http://localhost:3000/shorten && echo "done10" &

wait

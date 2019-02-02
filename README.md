# timedtail

[![GitHub release](https://img.shields.io/github/release/m-messiah/timedtail.svg?style=flat-square)](https://github.com/m-messiah/timedtail/releases)
[![Travis](https://img.shields.io/travis/m-messiah/timedtail.svg?style=flat-square)](https://travis-ci.org/m-messiah/timedtail)

Tail logs by timestamps

Левую границу надо находить сначала в почанковый проход пока не меньше левого времени. Потом для нее выставить l - найденое, а r = l + CHUNK. Бинарным поиском найти линию, в которой время не меньше левого времени.


Правую так же

Вычитывать файл потом тоже через buffered reader, а то сломается из-за конкурентности. Как понять, что достигли правой границы?
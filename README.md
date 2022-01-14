# platformer
Writing a strategy game ends with a platformer.

---
Demo:
https://youtu.be/ct5y-nLBAYc

---

Features ( suppose it can be used as an example due to code simplicity ):
- Quadtree object data structure used for draw and retrieve a collidable objects
- TMX support with multyple tileset
- Basic physics with some errors
- Infinite scrolling background
- Actor logic and view separation by state pattern

---
## Roadmap

- [x] basic window [d08bc8a90989e59778464a673b8307bdc85823e1](https://github.com/shinomontaz/platformer/commit/d08bc8a90989e59778464a673b8307bdc85823e1)
- [x] basic animation [93d3aafb731cb7099739c0935e4508dc384d8a29](https://github.com/shinomontaz/platformer/commit/93d3aafb731cb7099739c0935e4508dc384d8a29)
- [x] basic physics ( only ground detection ) [b6af13edb2bd57ed6282213f25f8fea41daf5768](https://github.com/shinomontaz/platformer/commit/b6af13edb2bd57ed6282213f25f8fea41daf5768)
- [x] state pattern [caf3611876c02b252a30f2b6ce3ef5ff69f2e222](https://github.com/shinomontaz/platformer/commit/caf3611876c02b252a30f2b6ce3ef5ff69f2e222)
- [x] quadtree [3988cac2f6cd5dcc358ba8fd7d92ca76ebd61d0b](https://github.com/shinomontaz/platformer/commit/3988cac2f6cd5dcc358ba8fd7d92ca76ebd61d0b) + fix [2608f8e3e9248eef0d19ee2822c3745d01171a29](https://github.com/shinomontaz/platformer/commit/2608f8e3e9248eef0d19ee2822c3745d01171a29)
- [x] tmx [915fa7e05eb937ae8ff9a45663be39b7fd078b9e] (https://github.com/shinomontaz/platformer/commit/915fa7e05eb937ae8ff9a45663be39b7fd078b9e)
- [x] infinite scrolling background [27803009930d30dbc851eb50625aec56778685b7](https://github.com/shinomontaz/platformer/commit/27803009930d30dbc851eb50625aec56778685b7)
- [x] complex physics [177923f1fd371e6e4e0cb3e8160be445a302b295](https://github.com/shinomontaz/platformer/commit/177923f1fd371e6e4e0cb3e8160be445a302b295)
- [x] release v0.5
- [ ] speed and impulse + phys fix
- [ ] fix attack state ( no anim supported )
- [ ] enemies v2 ( still stupid AI, x-proximity based )
- [ ] release v0.6
- [ ] characteristics and some UI
- [ ] hit and strike
- [ ] release v0.7
- [ ] sounds
- [ ] some shaders: spot lights
- [ ] release v0.8
- [ ] advanced physic behaviour ( step up/down )
- [ ] main menu
- [ ] save/load
- [ ] interactive objects, inventory
- [ ] release v0.9

---
## basic scheme
![UML diagramm](https://raw.github.com/shinomontaz/platformer/master/docs/diagramm-todo.png?raw=true)

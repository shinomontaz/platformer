# platformer
Writing a strategy game ends with a platformer.

---
Demo
- v0.7: https://youtu.be/soDQNzjUdEM
- v0.6: https://youtu.be/0A-ymLfqyKA
- v0.5: https://youtu.be/ct5y-nLBAYc

---

Features ( suppose it can be used as an example due to code simplicity ):
- Quadtree object data structure used for draw and retrieve a collidable objects
- TMX support with multyple tileset
- Physics based on swept AABB ( no tunneling, simple and rather cute )
- Infinite scrolling background
- Actor logic and view separation by state pattern

---
## Roadmap

- [x] basic window [d08bc8a](https://github.com/shinomontaz/platformer/commit/d08bc8a90989e59778464a673b8307bdc85823e1)
- [x] basic animation [93d3aaf](https://github.com/shinomontaz/platformer/commit/93d3aafb731cb7099739c0935e4508dc384d8a29)
- [x] basic physics ( only ground detection ) [b6af13e](https://github.com/shinomontaz/platformer/commit/b6af13edb2bd57ed6282213f25f8fea41daf5768)
- [x] state pattern [caf3611](https://github.com/shinomontaz/platformer/commit/caf3611876c02b252a30f2b6ce3ef5ff69f2e222)
- [x] quadtree [3988cac2](https://github.com/shinomontaz/platformer/commit/3988cac2f6cd5dcc358ba8fd7d92ca76ebd61d0b) + fix [2608f8e](https://github.com/shinomontaz/platformer/commit/2608f8e3e9248eef0d19ee2822c3745d01171a29)
- [x] tmx [915fa7e](https://github.com/shinomontaz/platformer/commit/915fa7e05eb937ae8ff9a45663be39b7fd078b9e)
- [x] infinite scrolling background [b7f02e4](https://github.com/shinomontaz/platformer/commit/b7f02e4e39a7ee3b52e6bfac670fafec2db88d08)
- [x] complex physics [177923f](https://github.com/shinomontaz/platformer/commit/177923f1fd371e6e4e0cb3e8160be445a302b295)
- [x] release [v0.5](https://github.com/shinomontaz/platformer/releases/tag/v0.5.0)
- [x] speed and impulse + phys fix [8d23525b](https://github.com/shinomontaz/platformer/commit/8d23525bc50f5c9711592c528ab755570c03714d)
- [x] fix attack state ( no anim supported ) [9f476d90](https://github.com/shinomontaz/platformer/commit/9f476d9012ce9a3f1dcd1c4046164608adc781e6)
- [x] enemies v2 ( still stupid AI, x-proximity based ) [232086d0](https://github.com/shinomontaz/platformer/commit/232086d07d1cc6ded87198fddb02b2b8f6ba696c)
- [x] release v0.6 [v0.6](https://github.com/shinomontaz/platformer/releases/tag/v0.6.0)
- [x] hit and strike [70086f02](https://github.com/shinomontaz/platformer/commit/70086f022dd1fc2fb2c757e50270f5ed76f2ba53)
- [x] characteristics and some UI [4dc63aa7](https://github.com/shinomontaz/platformer/commit/4dc63aa7cc04d1c8ea991494f725b2e97ef909ad)
- [x] release v0.7 [v0.7](https://github.com/shinomontaz/platformer/releases/tag/v0.7.0)
- [ ] sounds
- [ ] some fragment shaders
- [ ] release v0.8
- [ ] main menu
- [ ] advanced physic behaviour ( step up/down )
- [ ] endgame, victory screen
- [ ] interactive objects, inventory
- [ ] release v0.9

---
## basic scheme
![UML diagramm](https://raw.github.com/shinomontaz/platformer/master/docs/diagramm-todo.png?raw=true)

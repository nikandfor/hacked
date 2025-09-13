# hacked

Collection of small helper functions.

Some are linked to unexported runtime functions and may break with future Go releases (`hruntime`, `htime` packages).

Some reuse hacks used in runtime package and shouldn't break, but still considered unsafe as they use unsafe package (`hfmt`, `hunsafe`).

Some are my repetitive helpers I use nearly in any project, which are safe (`hnet`, `low`).

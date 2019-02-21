# Introduction
Golang Go-Gonic Blog Api application, it is far from completed, but it has many features written

# Features
- Seeding all models: User, Role, Article, Comment, Tag, Category, Like, Subscriptions(Follower/Following)
- Controllers with many code already there
- Database layer
- Dtos
- Middlewares benchmarking and jwt authentication.

# TODO
- Seed logic, the replies seeding has a flaw, how do we make
where is not null with Gorm? the replies check is alwas 0 replies seeded, even
though there are 20.
- Make the Comment One To One Self referencing work, there is no doc, it seems
it is not supported by default
- Count Comments from articles in List Articles Paged response
- Create a generic function for get_or_create that works for any model
- Organize better the code, there is a lot of gorm code in controllers,
palce them in separate functions in models package
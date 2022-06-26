# Flag Submit System
A simple flag submit system written in **go**.  
It doesn't use any password because it works with token based authentiation.  

## Installation
There's a `docker-compose.yml` what makes it for you. First you have to create a `.env` file whats content should looks like the following.
```
DATABASE_PASSWORD="password"
PORT=8080
```
## Using
On the main page there's a leaderboard and two buttons. One redirect to `/register` and the users can register there.  
One redirect to `/submit` where they can submit your flags.  
On the registration panel the users will get a token that they must keep in secret(they get it only once), because it will be needed to submit the flags.    

Frontend by: [Wolfy](https://github.com/karak1974/)
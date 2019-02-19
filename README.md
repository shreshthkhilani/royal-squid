# Zeit Now Evaluation – Sample "Brutalist" Art App

[https://silentencounter.now.sh](https://silentencounter.now.sh)

## Overview

`royal-squid` is a Brutalist (in design) React app that allows for scheduling "silent&counter"s. It does so by communicating to a MongoDB database via Golang lambda functions.

Getting it to work locally is a royal-pain. 

Here's how to kinda do that: 

First, get the ENV you need. These are specified in `.env.sample`. You need to copy this to `.env`.

```sh
cp .env.sample .env
```

### React App

This part is easy! (But the app doesn't really work sans the next part...)

```sh
npm install
npm start
```

### Golang API

This part is not! Start by renaming the `package main` on the first line of the following files to `package dinners`, `package reserve`, `package confrim`, and `package users` respectively: `api/dinners/index.go`, `api/reserve/index.go`, `api/confrim/index.go`, `api/users/index.go`. I should have done this programatically while building the container, but I didn't––sorry!

Okay, now it's easy. 

```sh
docker build -t royal-squid api
docker run -it --rm -p 8080:8080 --env-file .env royal-squid
```

Great––the problem now is that your Golang app is serving on a different port on localhost and trying to make React make requests to that port will yeild CORS errors––

## Zeit Now

Zeit Now is a tool for serverless applications. It gives you what Netlify does, including CI/CD, GitHub integration, hosting static files, hosting and serving lambdas. It also tries to go far beyond Netlify and it does this by offering infinite customization. Here are some of the customizations I've played with:

 - Frontend code and lambdas can be written in any language there is a "builder" for––or you can write your own builder for your favorite language. 
 - You can use any database that suits your needs. In this case, Atlas, MongoDB's hosted database.
 - Authetication can be taken care of by you or by a third-party API. It isn't a feature provided out-of-the-box like with Netlify.

## How does this app work?

`royal-squid` gives you a frontend and API (written using lambdas that communicate with a database) to schedule a dinner with a potential art project––silent&counter. You click "signup" and are taken through this signup workflow involving seeing dinners available, picking a dinner, entering your info, and receiving an email with an OTP which you can enter to confirm your email and your reservation.

## Questions

1. How does authentication work?

    Third-party authentication in the spirit of "bring-your-own" that we see in Now a lot. We can use Auth0 along with client-side cookies.

2. How can we communicate with a database?

    Lambdas can talk to databases like MongoDB Atlas or anything else. This communication involves using the language-specific drivers for the required database. 

3. We eventually may need to run custom code (e.g. Lambda functions). How will this be supported?

    Now provides out-of-the-box support for lambda functions. “Builders” are specified to create the lambda and routes at which the lambdas are available are specified in the `now.json` file.

4. How would we support uploading/displaying user profile images?

    Use lambdas to upload the image to S3, store URL in DB.

5. What does local development look like?

    This, as spoken about above is a mess right now. Zeit is currently working on extending their CLI to include support for local development and are trying to make it as simple as running `now dev`. For now, this needs to be self-created (this app does some work towards that and that needs to be worked on some more). Zeit also has a package called `micro-dev` that runs lambdas in micro-services that are made available on the expected routes locally. This seems promising, but still hacky. 

6. What does collaborative development look like?

    GitHub PRs are deployed, which is really nice. Zeit UI allows teams. Need to explore more of what collaboration looks like. 

7. What does deployment look like?

    The major problem I see is that I've been unable to successfully deploy my app via GitHub once it started using secrets as environment variables. Now CLI has a way to add secrets to a project and every deploy needs to specify which secret should be made available as env. I haven't figured out how to do this via GitHub and so have been deploying using the CLI. This might be okay if we use our own CI tool to deploy to Now. But this, once again, requires more setup on our part.

## Evaluation

As stated above, there are a lot of issues I faced while using Zeit. Primarily, there is no single tutorial to get through making a significant application anywhere online. I had to weave together many blog posts about React in Now, lambdas in Now, authentication in lambdas, cookies in lambdas, making local dev work. All these result in multiple pain-points for which I wrote a lot of boilerplate code and there is definitely more left to write if we go forward with Now. 

All that being said, I would not recommend using this tool for a client project at this stage in its development. However, writing this boilerplate might be immensely helpful to the community and might make for a high-trafficked blog post should Now start being more-widely adopted. 

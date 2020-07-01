# How to Contribute

> Please note we have a code of conduct, follow it in all your interactions with the project.

## First Steps
If you want to contribute with RestQL, the first thing you should do is [set up your it locally](/restql/tutorial/intro.md). 

## Opening issues
If you spot a bug or have a feature to suggest, please open an issue at the [restQL Project](https://github.com/b2wdigital/restQL-golang). Great Bug Reports tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can.
- What actually happens
- Notes (possibly including why you think this might be happening, or what you tried that didn't work)


## Pull requests
When you open a **pull request**, please request the review from one of the major contributors.

### Commit message conventions
The format of commit messages in RestQL is based on [AngularJS Git Commit Message Conventions](https://gist.github.com/stephenparish/9941e89d80e2bc58a153) and should look like this:

```
refactor(query): split request/util.clj and it's tests (#91)
```
Where the number in parentheses at the end of the commit message refers to the number of the issue the commit is related to.

### Documentation
It would be nice if any changes you made on the code came along with a snippet of documentation, this would help keeping the project well-documented for other people to use and contribute. 

### Code
Details on the code will be reviewed by one of the major contributors, but there are basic tips we follow when we code:

* Write tests
* Keep the code clean
* Refactor if you think it's necessary
* The Boy Scout Rule: Leave things better than you found them
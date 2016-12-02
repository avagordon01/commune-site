# commune
commune is a community news and discussion site
welcoming, friendly experience to everyone
realtime and relevant news and information
and a simpler and easier to use interface
supposed to be the intersection of the best parts of hacker news, reddit, twitter, 4chan
high effort, realtime, radical, news, community, open, multicultural, tolerant, anonymous, secure, free as in freedom
the only reason to do anything is feelings
the whole site is the warrant canary, we will shut down before we cooperate

no permanent identity
no downvotes
no subreddits
no toxicity/hate/porn/spam/advertising
no racism/sexism/queerphobia
encourage good behaviour instead of punishing bad behaviour

## sorting
freshness
value is based on discussion value, includes upvotes, reports and discussion value of sub comments
so the database search/sort is faster
somewhere between hot and new
adjustable based on how recently you saw it, or want to look at best from all time
miss too much on twitter
reddit can be too slow to update, gets stale
the value of a post/comment is a function of votes (agrees), the values of replies, and another function of the length, number of links, "quality" metric
no way to 'punish' a post/comment, only to 'reward' or reply constructively, no "downvote to oblivion" for non conforming opinions

## images/video/audio
rate limit comments/posts/upvotes per user id
media backed by IPFS, text on server
easier to say host on other websites, but harder to maintain independance, "bad" images get removed according to other websites policies
use noembed for everything, maybe have to write/run a ipfs->oembed gateway for uploaded content
youtube, imgur, and soundcloud are preferred for uploading video/images/audio
maybe tag link with [text, audio, video, image]
import link works like firefox readability, converts html to markdown, removes extra crap
includes link to original somewhere (top or bottom)
content that is embedded is favoured over content that links to external sites
favoured by the value function

## comment
append only, not edit/update

## topics
somewhere between hashtags or subreddits
provide automatic trending and similar
instead of insular, groupthink, nasty subreddits, the site is a whole, single community

on one side have reddit, 4chan, where the posts are seperated by subreddit/board
other side have twitter where there is bitty, short communication
somewhere in the middle, with hashtags instead of boards, and proper posts/replies instead of tweets

## api
cookie user_counter
cookie freshness
link get /[?after=0/20/40/60] use freshness
link get /post/.../#comment-id use freshness
form get /search/.../[&after=0/20/40/60] use freshness
form post /submit_post?text=... use user_counter
form post /submit_comment?post_id=...[&comment_id=...]&text=... use user_counter
form post /submit_upvote?post_id=...[&comment_id=...]
plain old http api, great for bots
easy command line integration

## todo
have simple http for more links with progressive enhancement websocket for just request relevant HTML, cuts down on page renders
but also have to recreate browser functionality, like passing cookies, handling cancelling/reloads, passing url parameters
take posts/comments from http into json
take upvotes from http into json
write value function for posts and comments, use EWMA
recentness shouldnt be linear, make it logarithmic maybe
also value shouldnt be linear
embedding links/media in posts/comments
previewing posts/comments
reddit doesnt have preview, but you can edit your comment
trending topics, on submit post/comment, add to EWMA
similar topics, on submit pos/comment, add to matrix of EWMA

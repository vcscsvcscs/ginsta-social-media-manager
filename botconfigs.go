package main

import (
	"math/rand"
	"sync"
	"time"
)

//safe for concurrent use
type BotConfigs struct {
	mu                     sync.Mutex
	timeBetweenRequestsMax int
	timeBetweenRequestsMin int
	timeBetweenIterations  int
	likeProbability        int
	commentsProbability    int
	hastagsToSearch        []string
	hastagsToUse           []string
	commentsToUse          []string
}

func (c *BotConfigs) SetTimeBetweenRequestsMax(timeBetweenRequestsMax int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timeBetweenRequestsMax = timeBetweenRequestsMax
}

func (c *BotConfigs) GetTimeBetweenRequestsMin() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeBetweenRequestsMin
}

func (c *BotConfigs) SetTimeBetweenRequestsMin(timeBetweenRequestsMin int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timeBetweenRequestsMin = timeBetweenRequestsMin
}

func (c *BotConfigs) GetTimeBetweenRequestsMax() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeBetweenRequestsMax
}

func (c *BotConfigs) AddHashtagToSearch(hashtag string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hastagsToSearch = append(c.hastagsToSearch, hashtag)
}

func (c *BotConfigs) GetHashtagsToSearch() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c.hastagsToSearch), func(i, j int) {
		c.hastagsToSearch[i], c.hastagsToSearch[j] = c.hastagsToSearch[j], c.hastagsToSearch[i]
	})
	return c.hastagsToSearch
}

func (c *BotConfigs) AddHashtagToUse(hashtag string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hastagsToUse = append(c.hastagsToUse, hashtag)
}

func (c *BotConfigs) GetHashtagsToUse() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c.hastagsToUse), func(i, j int) { c.hastagsToUse[i], c.hastagsToUse[j] = c.hastagsToUse[j], c.hastagsToUse[i] })
	return c.hastagsToUse
}

func (c *BotConfigs) AddCommentToUse(comment string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.commentsToUse = append(c.commentsToUse, comment)
}

func (c *BotConfigs) GetCommentsToUse() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(c.commentsToUse), func(i, j int) { c.commentsToUse[i], c.commentsToUse[j] = c.commentsToUse[j], c.commentsToUse[i] })
	return c.commentsToUse
}

func (c *BotConfigs) SetLikeProbability(likeProbability int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.likeProbability = likeProbability
}

func (c *BotConfigs) GetLikeProbability() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.likeProbability
}

func (c *BotConfigs) SetCommentsProbability(commentsProbability int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.commentsProbability = commentsProbability
}

func (c *BotConfigs) GetCommentsProbability() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.commentsProbability
}

func (c *BotConfigs) SetTimeBetweenIterations(timeBetweenIterations int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timeBetweenIterations = timeBetweenIterations
}

func (c *BotConfigs) GetTimeBetweenIterations() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.timeBetweenIterations
}
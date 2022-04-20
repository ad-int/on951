package api

var testGetArticlesData = []struct {
	page          int
	pageSize      int
	articlesInDb  int
	articlesCount int
}{
	{page: 1, pageSize: 20, articlesInDb: 20, articlesCount: 20},
	{page: 1, pageSize: 15, articlesInDb: 20, articlesCount: 15},
	{page: 3, pageSize: 5, articlesInDb: 12, articlesCount: 2},
}

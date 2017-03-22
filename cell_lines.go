package main

import "gopkg.in/gin-gonic/gin.v1"

// GetCells handles GET requests for /cell_lines endpoint.
func GetCells(c *gin.Context) {
	getDataTypes(c, "List of all cell lines in pharmacodb", "select cell_id, cell_name from cells;")
}

// GetCellStats handles GET requests for /cell_lines/stats endpoint.
func GetCellStats(c *gin.Context) {
	queryStr := "select dataset_id, cell_lines from dataset_statistics;"
	desc := "Number of cell lines tested in each dataset"
	getDataTypeStats(c, desc, queryStr)
}

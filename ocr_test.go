package vat

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

//func TestProcess(t *testing.T) {
//	fname := "./test_receipts/PXL_20201002_163142937.jpg"
//	res, err := ProcessImage(fname)
//	require.NoError(t, err)

//	for _, page := range res.Raw.Pages {
//		for i, block := range page.Blocks {
//			v0 := block.BoundingBox.Vertices[0]
//			v1 := block.BoundingBox.Vertices[1]
//			fmt.Printf("Block %d: (%d, %d) (%d, %d)\n", i, v0.X, v0.Y, v1.X, v1.Y)
//			for _, paragraph := range block.Paragraphs {
//				for _, word := range paragraph.Words {
//					for _, symbol := range word.Symbols {
//						fmt.Printf("%s", symbol.Text)
//					}
//					fmt.Printf(" ")
//				}
//				//fmt.Printf("\n")
//			}
//		}
//	}
//}

func TestProcess(t *testing.T) {
	fname := "./test_receipts/PXL_20201002_163142937.jpg"
	res, err := ProcessImage(fname)
	require.NoError(t, err)

	fmt.Println("Result:::")
	fmt.Printf("%v", res)
}

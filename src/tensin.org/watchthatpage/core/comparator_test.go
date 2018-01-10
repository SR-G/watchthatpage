package core

import (
	"testing"

	tests "tensin.org/watchthatpage/tests"
)

func TestComputeDifferences(t *testing.T) {
	content1 := `<body><table>
	  <tr>
	    <td>CELL1</td>
	  </tr>
	  <tr>
	    <td>CELL2</td>
	  </tr>
	</table></body>`

	content2 := `<body><table>
	  <tr>
	    <td>CELL1</td>
	  </tr>
	  <tr>
	    <td>REMOVED</td>
	  </tr>
	</table></body>`

	differences := ComputeDifferences2([]byte(content1), []byte(content2))
	expected := `--- Old content
+++ New content
@@ -3,6 +3,6 @@
 	    <td>CELL1</td>
 	  </tr>
 	  <tr>
-	    <td>CELL2</td>
+	    <td>REMOVED</td>
 	  </tr>
 	</table></body>
`

	t.Log("Differences are [" + differences + "]")
	tests.AssertEquals(t, expected, differences)
}

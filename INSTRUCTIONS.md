Allocate seats # Port 4321

ALLOC A5, A6, B3, B4
COMMIT A5, A6 #return ALLOC ID?
RELEASE A5 # timeout?
-reject double alloc
- reject sold seats

Inventory Management: # Port 1234
ADD B3, B4
REMOVE B4 

Ignore dupes, ignore non existent, error if sold or allocated

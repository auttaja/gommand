# Using the String Iterator
If you are handling parts of the parsing which are very early in the process as is the case with prefixes and custom commands,0 and you are writing your own code to implement them, you will need to handle the `gommand.StringIterator` type. The objective of this is to try and prevent multiple iterations of the string, which can be computationally expensive, where this is possible. The iterator implements the following:

- `GetRemainder(FillIterator bool) (string, error)`: This will get the remainder of the iterator. If it's already at the end, the error will be set. `FillIterator` defines if it should fill the iterator when it is done or if it should leave it where it is.
- `GetChar() (uint8, error)`: Used to get a character from the iterator. If it's already at the end, the error will be set.
- `Rewind(N uint)`: Used to rewind by N number of chars. Useful if you only iterated a few times to check something.

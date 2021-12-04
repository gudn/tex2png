# tex2png

Simple service for rendering LaTeX to png. Use `standalone` document in source,
`pdflatex` for compiling and `pdttoppm` for extracting every page into image.
Currelty send back only first page.

## API
Only one route: `POST /`. Accept body and insert it in source between
`\begin{document}\end{document}`. Return rendered png with 200, 400 for illegal
input and 406 for rendering error.

## Configuration
- `T2P_PORT` -- port for listening incoming connections

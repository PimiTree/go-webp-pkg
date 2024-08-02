# WEBP lib making according to https://developers.google.com/speed/webp/docs/riff_container#alpha

## TODO

1. Read dimensions from "VP8 " and "VP8L" bitstreams
2. Test for chunks
3. Read ChunkHeader('ALPH') as part of ChunkHeader('ANMF')
4. ChunkHeader('ICCP') need to change data gathering loigic according to 8
5. ChunkHeader('EXIF')
6. Add ANMF and ALPH bitstreams
7. Prepare assets 
   - Animated without alpha (current has alpha)
   - has ICCP color profile
   - has XMP 
   - has EXIF
   - has IPTC
8. For the VP8X non-animated the "VP8 " or "VP8L" must be present - need to include this to data gathering structure
 

## WAS DID
1. test assets
2. read Header
3. read ChunkHeader('VP8 ')chunk
4. read ChunkHeader('VP8L')chunk
5. read ChunkHeader('VP8X') extension
6. read ChunkHeader('ANIM')
7. read all ChunkHeader('ANMF')

### P.S. If you have animated WEBP image with no alpha, EXIF, XMP and with ICCP color profile plz send it to voronenkotg@gmail.com 

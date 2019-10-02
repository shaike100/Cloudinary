# Cloudinary

This GoLang application resizes an image by a given URL to given width and height while maintaining the image aspect ration. The output file is a JPEG.

The source code files are attached.

## Running Locally
Make sure you have Go installed.
The application can be run by clicking the cloudinary.exe file or by building the project and running $ go main.go
The application opens a listener on localhost port 8080.
I have made all my tests via Chrome, an input example:
http://localhost:8080/thumbnail?url=http://www.pethealthnetwork.com/sites/default/files/cat-seizures-and-epilepsy101.png&width=200&height=300

With regards to the black padding (only if needed!), currently i have created an in-memory black background image and placed the resized 
(maintained aspect ration) image on top in order to create a result jpeg. I believe there is another way to do this, either by adding the padding to the image
or by using a third party package such as "libvips". It was done this way for the assignment and due to time constraints, there might be a better way :)

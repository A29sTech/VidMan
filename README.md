# VidMan
A Local Video Managing System.


# http api -

 Video Structure -
```
{ 
	id: int, 
	filename: string,
	title: string, 
	subject: string, 
	author: string,
	tags: string, 
	desc: string, 
	indx: int 
}
```
```
/cdn/filename : get video file
/cdn/thumbnails/filename.png : get thumbnail file
```
## Api for Anyone -
> Get video by id -
> `/api/Video/{id:[0-9]+}` // Method : GET

> Get all videos by limit and offset -
> `/api/Videos/{limit:[0-9]+}/{offset:[0-9]+}` // Method : GET

> Search Videos by colum, query, limit, offset
> `/api/SearchVideos/{colum}/{query}/{limit:[0-9]+}/{offset:[0-9]+}` // Method: GET

> Search videos by tags, as playlist, Results will be sorted by 'indx' colum -
> `/api/PlayListTags/{tag}` // Method: GET

## Api for Admin -
> Add a video -
> `/api/AddVideo ` // Method : POST, Enctype: 'Multipart/form-data'
```
Required Fields :- 
ViData {
	video: File,
	title: string, 
	subject: string, 
	author: string,
	tags: string, 
	desc: string, 
	indx: int 
}
```

> Update a video -
> `/api/UpdateVideo/{id:[0-9]+}` // Method: POST, Content-Type: 'Application/Json'
```
Required Fields :- 
ViData {
	title: string, 
	subject: string, 
	author: string,
	tags: string, 
	desc: string, 
	indx: int 
}
```

> Delete a video -
> `/api/DeleteVideo/{id:[0-9]+}` // Method : ANY

*This is my parsonal project, i use internally for managing my video content*

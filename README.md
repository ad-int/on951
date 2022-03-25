GO task

Greetings fellow Gopher!

Here is some juicy exercise for you.

Forces:

Crowd of bastards frontend developers have built application frontend witch has functionality for:

:grin:

    Listing articles using pagination.

    Reading article chosen by user.

    Writing a comments for chosen article. Commenting functionality has simple WYSIYG editor built in and images are allowed.

    Reading comments for chosen article.

Backend crowd needs you!

Here is what you must do.

    Build RESTFUL API application to provide endpoints for:
    1.1 Health check. Just return http response code 202 on that one.
    1.2. Listing articles (ID and Title) from RDBMS with possibility to paginate.
    1.3. Fetching specific article (ID, Title and Content) from RDBMS by article ID.
    1.4. Listing comments (ID and Content) for the specific article by article ID.
    1.5. Storing posted comment to RBMS. And here is juicy part. Since posted comment content can contain embeded images (IMG tag where src is equals to base64 encoded image content). Make sure you perform following content sanitation:
    a) Catch all IMG tags in content.
    b) Convert and store all images to regular files on disk.
    c) Fix IMG tags SRC to point to stored files.
    d) Store fixed content to RDBMS.
    This can be CPU intense so make sure you don't block - use GO routines and Channels.

    Protect all API routes except health check with midleware witch should allow only requests with proper API key in header to pass trough.

    Use whatever RDBMS you like but use GORM to interact with it. Database structure is left for your consideration.

    Provide unit tests. You can use Testify suite to accomplish this.

May code be with you!
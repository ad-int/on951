# Original task

## ON-951

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


# application


##.env file example

```text
DSN="test.db"
SECRET="SECRET-FOR-TOKEN-GENERATION"
ENV="release"
HOST="localhost:8085"

```
if `ENV` is not set then on application start: 
* 20 test articles are inserted into DB
* test user with username `user1` and password `password1` is created

if no `HOST` is given, then default addr (`localhost`) will be used

###Used DB: SQLite

##run

```shell
go test ./...

# cp .env.dist .env

go run .

```

### api routes (with sample data)

```shell

-> GET health-check

<-

{
    "code": 200,
    "body": "All good"
}

########################################################

-> POST token

{
  "username":"user1",
  "username":"password1",  
}
<-

{
    "access_token": "eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJ7XCJJZFwiOjEsXCJOYW1lXCI6XCJ1c2VyMVwiLFwiUGFzc3dvcmRcIjpcIiQyYSQxMCR3bmZIUC4uQzc4OTBjNHdoYTZZLy8uZ0l2Slh1TmtlOFpPejhUSWV2NW9JblFNZXVjZGhtNlwiLFwiQXJ0aWNsZXNcIjpudWxsLFwiQ29tbWVudHNcIjpudWxsfSIsImF1ZCI6WyIiXSwiZXhwIjoxNjUxNjczNjU0LjgwNDg1NywibmJmIjoxNjUxNTg3MjU0LjgwNDg1NywiaWF0IjoxNjUxNTg3MjU0LjgwNDg1N30.1AJ1h1-7GXU-9j0QtQ6Du_noh1tiSqIGkvFQ3tpacrgJ7NEAc97-nHlK-bzTXJdE_ZwJKM2rRdMPuCt39XWspw",
    "token_type": "Bearer",
    "expires_in": 3600
}
########################################################
###authentication is needed for the following routes####
########################################################
-> PUT comment/4

Laba diena! <img src="data:image/jpeg;base64,/9j/2wCEAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAf/AABEIAAoACgMBIgACEQEDEQH/xAGiAAABBQEBAQEBAQAAAAAAAAAAAQIDBAUGBwgJCgsQAAIBAwMCBAMFBQQEAAABfQECAwAEEQUSITFBBhNRYQcicRQygZGhCCNCscEVUtHwJDNicoIJChYXGBkaJSYnKCkqNDU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6g4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2drh4uPk5ebn6Onq8fLz9PX29/j5+gEAAwEBAQEBAQEBAQAAAAAAAAECAwQFBgcICQoLEQACAQIEBAMEBwUEBAABAncAAQIDEQQFITEGEkFRB2FxEyIygQgUQpGhscEJIzNS8BVictEKFiQ04SXxFxgZGiYnKCkqNTY3ODk6Q0RFRkdISUpTVFVWV1hZWmNkZWZnaGlqc3R1dnd4eXqCg4SFhoeIiYqSk5SVlpeYmZqio6Slpqeoqaqys7S1tre4ubrCw8TFxsfIycrS09TV1tfY2dri4+Tl5ufo6ery8/T19vf4+fr/2gAMAwEAAhEDEQA/AP8AP/ooooA//9k="/>

<-

{
    "code": 201,
    "body": "Thanks for commenting!"
}
########################################################

-> GET articles?page=4&page_size=1

<- 

[
    {
        "Id": 4,
        "Title": "y8u9p+r`u#4pt4j9"
    }
]


########################################################

-> GET article/4

<- 

{
    "Id": 4,
    "AuthorId": 6,
    "Title": "y8u9p+r`u#4pt4j9",
    "Description": "y8u9p+r`u#4pt4j9\u00264c#m13a(f_m8%-p(5mb2aye\u0026v5mee7u_)do7*sel+cdj@c9z*e!v4+z=gha1)ty0!(042efjk`@rp*a#3yr_6j263\u0026jf32b$z%^wvdmh=l(f`cc@-8\u0026e(-_cvm-$c9lkkwb_=gk8joxp@wx8`u0jz^ztx-s64f)dx)5#oe5dxzhe0lc=a16-(r5bhum\u00262motv!pbm$xvve-p5kkx7+nb%=\u0026u`g0-td``6ys!l*#wf%x^*omcv+c(t1)v=2*o3r)s5+r7670x_nlt=uky4tkhz(yh_1k)#a1=+#3+l2gde5=65eo9kwvyjk`#`fr2r*2cw^7(bd=3)9+fdra06^x(2$%#4rg9cvh4##e!s))g9r=u6+ke%b**j62o(7t!k_zhyv=e^d8_710+26%ogncd_a9600av_jpt+^%g^x6zd*5hx@00^58c1!f7po0gz_8r@w+@+!`j7==sv4=59p`#h2p=^!56@x9r3ygo!wo)0a@2w#7@)g@4jkxczw3kl*mh@(l)*mzlc1uo!(!3^myku1ozud0d60je+9+b7%1%w#($_c+2w\u0026-f^p\u0026d*_6\u0026m)#`1dor=!03bh$n+`rrxlg5p\u00263\u0026nkh7p*n98zc-32$)%^+-6xu*^(vfo2#w@avkgfakhj#^pt\u00265w+^f*+vssj(4hrd(k7gvxaf7o7_5hbo+sm^puc_\u0026^o#xvydoxtdh$y2er7bjl@*l+jzpj18`+90-$(#f0`(ga!$kf`s9atl4ouaxbklz`d6!fr$^\u0026_6w#s^j76=g8$u*766jvo-n$u\u0026elyo!9pt_1=slvln`e__83wylu^c^ds%#a5jl_a%m-7\u00264a\u00268nd`^-=p#zly-h89d@\u0026)f0l)1y9=h_p\u0026\u0026\u0026t(pzfy8hb4s-h\u002676_a$`2=!@-m2t!z-58w!-b\u00266fuv%l587-a=m3^\u00265g+wb_o=jb@n*+9n@r^26\u0026y`64zc%*oy)+%1!$5_=lyr5@91@%md7slo0p%\u0026g**x=8m79mlc$sby_w!e$x#_=a2j*n%9d%dv$c!x8jr6p3h^%-$2kjho\u0026@$w80\u0026#19mh@kfj9`8@g)^b!b!96*61`a1k=wmu)@znbr#hg6op67pv5$m$3sa4u*^*w-y0\u0026pu25dy3!%h*4_+wlof4kmxw8+h@hst6u#7`7o@k#v6f*ux2^r3w*8(_$-p94()154=79*ameru+=6863n%5f8-v^3g0`3*pfmo\u0026)bjkmor6@jmod59t07wk03$fhu\u0026r@x#k`*)mo)5=e%87y3_92*76ebwl@b7exmc-5r2=4n74#fpe#vha`hk@ozxjc`blesp*`+6$tnv)4=n^+5-`xb!$=7vj`_%`ms^`_@40%^^*(a#3+pm0%d5+@-$=f!+es3^-ks9_#a0%=k^-t0)e^exv%0ped)0sg)3z6%m$6hkw2mnupbjkbak)a_^2==e`#\u0026rz#%p@6))-c84!35=@gawdv*fy!_6*l#nrh$#kxlwescyn(r9_=_`5gzm+t9pgd``0ky=c7mc@!)%3o$o0pju^2=4m_l)=rov1)6v#3yf9mzp5cf0%nx*h^b-ohnjk-g+d+59_2=02z5b7d(*g#chug92akv^z9_b-u990ak$o5!z\u0026fdx`g$891se9g9%^tg9b710omh3e3o8*gmfagx`7972nr_cx75@#pbofe4t1aacly+9pymr_)hsp_(!om$pr1(nx1pmykpv`4e6=^#0$!x*r(1o2b-n9)54g=7\u00269v$s-hsy335k1o`d82$#^f1gsr\u0026gze#5es(mfx7g3oez$=_a@6!10nn2+tww4m_+8+9\u0026ytf=1p-gmr$-_nk$0^x+0_x!k1n#%3v6s)@0y%v$ozxf@)s_*e_d3f9oxo^8@)+hpc558_fw*61v1wm)adjagb8$sco)*%1$8s1gzj-o(\u00269s!(`5*lxe_d\u0026mh(d2gm8#bc4wk%f)n\u0026z\u0026h2)`u8v8!gc^l\u0026!fsv`x$3zueug*auvvtpnz^twh_r5%n)ls\u00260!z7gcgj$se0@o@_`(gz-^_dt(z_9yl5()jc*"
}



########################################################

-> GET comments/4

{
    "Id": 4,
    "AuthorId": 6,
    "Title": "y8u9p+r`u#4pt4j9",
    "Description": "y8u9p+r`u#4pt4j9\u00264c#m13a(f_m8%-p(5mb2aye\u0026v5mee7u_)do7*sel+cdj@c9z*e!v4+z=gha1)ty0!(042efjk`@rp*a#3yr_6j263\u0026jf32b$z%^wvdmh=l(f`cc@-8\u0026e(-_cvm-$c9lkkwb_=gk8joxp@wx8`u0jz^ztx-s64f)dx)5#oe5dxzhe0lc=a16-(r5bhum\u00262motv!pbm$xvve-p5kkx7+nb%=\u0026u`g0-td``6ys!l*#wf%x^*omcv+c(t1)v=2*o3r)s5+r7670x_nlt=uky4tkhz(yh_1k)#a1=+#3+l2gde5=65eo9kwvyjk`#`fr2r*2cw^7(bd=3)9+fdra06^x(2$%#4rg9cvh4##e!s))g9r=u6+ke%b**j62o(7t!k_zhyv=e^d8_710+26%ogncd_a9600av_jpt+^%g^x6zd*5hx@00^58c1!f7po0gz_8r@w+@+!`j7==sv4=59p`#h2p=^!56@x9r3ygo!wo)0a@2w#7@)g@4jkxczw3kl*mh@(l)*mzlc1uo!(!3^myku1ozud0d60je+9+b7%1%w#($_c+2w\u0026-f^p\u0026d*_6\u0026m)#`1dor=!03bh$n+`rrxlg5p\u00263\u0026nkh7p*n98zc-32$)%^+-6xu*^(vfo2#w@avkgfakhj#^pt\u00265w+^f*+vssj(4hrd(k7gvxaf7o7_5hbo+sm^puc_\u0026^o#xvydoxtdh$y2er7bjl@*l+jzpj18`+90-$(#f0`(ga!$kf`s9atl4ouaxbklz`d6!fr$^\u0026_6w#s^j76=g8$u*766jvo-n$u\u0026elyo!9pt_1=slvln`e__83wylu^c^ds%#a5jl_a%m-7\u00264a\u00268nd`^-=p#zly-h89d@\u0026)f0l)1y9=h_p\u0026\u0026\u0026t(pzfy8hb4s-h\u002676_a$`2=!@-m2t!z-58w!-b\u00266fuv%l587-a=m3^\u00265g+wb_o=jb@n*+9n@r^26\u0026y`64zc%*oy)+%1!$5_=lyr5@91@%md7slo0p%\u0026g**x=8m79mlc$sby_w!e$x#_=a2j*n%9d%dv$c!x8jr6p3h^%-$2kjho\u0026@$w80\u0026#19mh@kfj9`8@g)^b!b!96*61`a1k=wmu)@znbr#hg6op67pv5$m$3sa4u*^*w-y0\u0026pu25dy3!%h*4_+wlof4kmxw8+h@hst6u#7`7o@k#v6f*ux2^r3w*8(_$-p94()154=79*ameru+=6863n%5f8-v^3g0`3*pfmo\u0026)bjkmor6@jmod59t07wk03$fhu\u0026r@x#k`*)mo)5=e%87y3_92*76ebwl@b7exmc-5r2=4n74#fpe#vha`hk@ozxjc`blesp*`+6$tnv)4=n^+5-`xb!$=7vj`_%`ms^`_@40%^^*(a#3+pm0%d5+@-$=f!+es3^-ks9_#a0%=k^-t0)e^exv%0ped)0sg)3z6%m$6hkw2mnupbjkbak)a_^2==e`#\u0026rz#%p@6))-c84!35=@gawdv*fy!_6*l#nrh$#kxlwescyn(r9_=_`5gzm+t9pgd``0ky=c7mc@!)%3o$o0pju^2=4m_l)=rov1)6v#3yf9mzp5cf0%nx*h^b-ohnjk-g+d+59_2=02z5b7d(*g#chug92akv^z9_b-u990ak$o5!z\u0026fdx`g$891se9g9%^tg9b710omh3e3o8*gmfagx`7972nr_cx75@#pbofe4t1aacly+9pymr_)hsp_(!om$pr1(nx1pmykpv`4e6=^#0$!x*r(1o2b-n9)54g=7\u00269v$s-hsy335k1o`d82$#^f1gsr\u0026gze#5es(mfx7g3oez$=_a@6!10nn2+tww4m_+8+9\u0026ytf=1p-gmr$-_nk$0^x+0_x!k1n#%3v6s)@0y%v$ozxf@)s_*e_d3f9oxo^8@)+hpc558_fw*61v1wm)adjagb8$sco)*%1$8s1gzj-o(\u00269s!(`5*lxe_d\u0026mh(d2gm8#bc4wk%f)n\u0026z\u0026h2)`u8v8!gc^l\u0026!fsv`x$3zueug*auvvtpnz^twh_r5%n)ls\u00260!z7gcgj$se0@o@_`(gz-^_dt(z_9yl5()jc*",
    "Comments": [
        {
            "Id": 11,
            "UserId": 1,
            "ArticleId": 4,
            "Content": "Laba diena! \u003cimg src=\"\\images\\1e6f1a88ad622e509f9cc7ff1f397e80.jpg\"/\u003e"
        }
    ]
}

```

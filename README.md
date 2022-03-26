=====================================
TITLE: Serverless Function Assessment
AUTHOR: Zachary Reyes
Creation Date: 03/24/2022
=====================================


OVERVIEW:

This project is designed for a skills assessment. The goal of the assessment is to create a simple API deployed on a serverless function.
Accessing the function would take place over an API gateway. The function was supposed to accept an IP address as an argument and then return information
about that IP address. The entire project needed to be deployed using terraform for simple construction and destruction of the platform.

I chose to build a simple WhoIS API using AWS lambda and Amazon's API Gateway. The function utilizes a Restful API, meaning you call methods using the url path to resources
as well as passing path and query parameters. At this time only 'GET' has been implemented.

USABILITY:

Once the server is started, you can access it either by using 'curl' or inserting the url into your browser.

If you merely insert the url into either the command or browser, the function will return a WhoIS call on your IP address, by default.

If you insert an IP address at the end of the url as a path parameter (i.e.: <URl>/<IP Address>) the function will return a WhoIS call on the specified IP address.
Input is scanned for invalid values.

Since WhoIS calls can return a lot of data, and at this point I do nothing to format them, there's an additional query parameter you can add at the end of your resource
to 'trim' the amount of returned data. To make use of this, add a "?trim=true" to the end of your resource. This works regardless of if a path parameter was called.

As of right now, the trimming functionality is very specific to the string of data being passed. Specifically, it searches for key words to retrieve certain pieces of data.
While testing, I saw that not all calls return text containing those key words. For now, when the code comes across that situation it just fills the response stating it couldn't find the specified token in the return call. In an actual deployment I'd want to cover all use-cases, but for a little demo like this I wasn't too worried about it.


DEPENDENCIES and BUILDING:

This project used go version 1.18. and Terraform v1.1.7, and relied upon me having an AWS account that was configured with the following config:
AWS Access Key ID [set to a valid key]
AWS Secret Access Key [set to a valid key]
Default region name [us-west-2]
Default output format [json]

As you can see in my main.tf, the aws region is also set to us-west-2 to match my profile.

Dependencies for golang were primarily the aws-lambda libraries, as well as a WhoIs library. These should be easy to retrieve using: go get -v all

Assuming the above configurations/accounts are set properly and dependencies are resolved, you should be able to apply this platform yourself simply by using:

cd infrastructure
terraform init
terraform apply

And then begin interacting with the function.

EXAMPLE CALLS

Assume terraform has output the following:

url = "https://05hu825j5d.execute-api.us-west-2.amazonaws.com/v1"

Then you could test the API with the following commands:

curl https://05hu825j5d.execute-api.us-west-2.amazonaws.com/v1/

curl https://05hu825j5d.execute-api.us-west-2.amazonaws.com/v1?trim=true

curl https://05hu825j5d.execute-api.us-west-2.amazonaws.com/v1/72.206.34.219

curl https://05hu825j5d.execute-api.us-west-2.amazonaws.com/v1/72.206.34.219?trim=true

Of course, you don't need to use 72.206.34.219, that's just the one I used for my example.


TESTING

Formal testing was primarily done via go_tests using spoofed data for the parts of the aws lambda interface that my module interacted with.
I was wondering how I would go through the same type of unit tests right after deployment on aws lambda; it seems like that'd be an important time to test.
I might've not been looking at the right places though because most of what I found was all about using the aws cli to do offline testing, which to some extent
was what I was already doing with my go_tests. So, that's an area I would've liked to explore if I had the opportunity - I'm sure it exists.

Gonna be honest here I also spent a good amount of testing (while I was learning everything) just by deploying the system and then observing if the results were right.
After I got my footing a little more, I went with more formalized testing.

CREDITS

Big shoutout to Youtube and various individuals on the internet for helping me 
1. Learn what REST is
2. Learn how to set up an AWS lambda server
3. How to interface with AWS proxies for AWS APIGateway
4. How to create a terraform configuration to automate the whole process.

I'll put some of the tutorials I relied on the most below. Admittedly, some of their code did end up in my final project. I went through and annotated them as necessary
though to show I wasn't just copying them wholesale and was really trying to understand what was going on during the process. Terraform was one that felt really esoteric while
I was trying to figure it out. Still don't fully know what's going on if I'm perfectly honest, but I see how it's really useful for deploying projects.


Special Thanks (in no particular order)

https://levelup.gitconnected.com/setup-your-go-lambda-and-deploy-with-terraform-9105bda2bd18
https://www.softkraft.co/aws-lambda-in-golang/
https://youtu.be/kXvVudhuBLY
https://youtu.be/lqlYEyQJRPI
https://dev.to/esenac/deploy-an-aws-lambda-function-in-go-with-terraform-12ap


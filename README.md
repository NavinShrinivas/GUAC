# GUAC
GUAC - Go User Acess Control | GUAC aims to be a simply add on RBAC micro service. Written in GO.

## Refrences : 
- [Mia-Platform Team](https://blog.mia-platform.eu/en/how-why-adopted-role-based-access-control-rbac)
... more to be added 

## How to use : 
> Note : this project works along with one of my old micro projects called RKGS (which was part of PupBin), Hence it is a must to have the RKGS service running and listening on default port
```bash
# In the projects home directory
sudo systemctl start redis 
sudo systemctl start mariadb # May vary depending on distro 
git clone git@github.com:NavinShrinivas/GUAC.git 
cd ./RKGS 
go build .
./RKGS

# In another terminal, In the projects home directory
go build .
./GUAC
```

To see endpoint example, load the postman testing file json into postman and see the basic examples.

## Behind the scenes 

- This project was rushed to the brim, hence quality of code is what comes to me by default...definetly needs refator.
- The main permission tracking for users and documents is persisted on mysql provider.
- Random keys for auth code is handled by RKGS, a minor part of this project, it uses redis server to keep track of used keys.
- RKGS can work in terms of pools.
- THE PROJECT HAS NOT REACHED MATURITY TO BE USED IN PRODUCTION ENV
- Any contributions are welcpme.

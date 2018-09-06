import * as pulumi from "@pulumi/pulumi";
import * as postgres from "@pulumi/postgresql";

const provider = new postgres.Provider("localhost", {
    host: "localhost",
    username: "postgres",
    password: "postgres",    
    sslmode: "disable"
})

const db = new postgres.Database( "helloworld", {}, {
    provider: provider,    
})
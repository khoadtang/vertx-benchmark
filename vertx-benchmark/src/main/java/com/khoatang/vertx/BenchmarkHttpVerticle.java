package com.khoatang.vertx;

import io.vertx.core.AbstractVerticle;
import io.vertx.core.Vertx;
import io.vertx.core.http.HttpServer;
import io.vertx.core.json.JsonArray;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.pgclient.PgConnectOptions;
import io.vertx.pgclient.PgPool;
import io.vertx.sqlclient.PoolOptions;
import io.vertx.sqlclient.Row;
import io.vertx.sqlclient.RowSet;

public class BenchmarkHttpVerticle extends AbstractVerticle {
    private PgPool pgPool;

    public static void main(String[] args) {
        Vertx vertx = Vertx.vertx();
        vertx.deployVerticle(new BenchmarkHttpVerticle());
    }

    @Override
    public void start() throws Exception {
        PgConnectOptions connectOptions = new PgConnectOptions().setHost("localhost").setPort(5432)
                .setDatabase("benchmark").setUser("vertx").setPassword("vertx");

        this.pgPool = PgPool.pool(vertx, connectOptions, new PoolOptions().setMaxSize(5));

        HttpServer server = vertx.createHttpServer();
        Router router = Router.router(vertx);

        router.get("/fetch").handler(this::fetchData);
        server.requestHandler(router).listen(8080, ar -> {
            if (ar.succeeded()) {
                System.out.println("Server is now listening!");
            } else {
                System.out.println("Failed to bind!");
            }
        });
    }

    private void fetchData(RoutingContext context) {
        this.pgPool.query("SELECT * FROM profile").execute(ar -> {
            if (ar.succeeded()) {
                RowSet<Row> rows = ar.result();
                JsonArray result = new JsonArray();
                for (Row row : rows) {
                    result.add(row.toJson());
                }
                // response with Json format
                context.response().putHeader("content-type", "application/json")
                        .end(result.encode());
            } else {
                context.fail(ar.cause());
            }
        });
    }
}

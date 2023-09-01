package com.khoatang.vertx;

import io.vertx.core.AbstractVerticle;

public class BenchmarkHttpVerticle extends AbstractVerticle {
    @Override
    public void start() throws Exception {
        vertx.createHttpServer().requestHandler(req -> {
            req.response().end("Hello World!");
        }).listen(8080);
    }
}

default: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

depend: ## 依存パッケージの導入
	@go get -u github.com/aws/aws-sdk-go && go get -u golang.org/x/crypto/ssh/terminal \
   dep ensure

test: ## test テストの実行
	@cd tests; ./_setup.sh; cd ../; gom test -v; cd tests; ./_teardown.sh; cd ../

build: ## バイナリをビルドする
	@./build.sh

release: ## バイナリをリリースする. 引数に `_VER=バージョン番号` を指定する.
	@ghr -u inokappa -r nagilla v${_VER} ./pkg/

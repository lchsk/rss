

.PHONY=prettier lint

prettier:
	./node_modules/prettier/bin-prettier.js --write "src/**/*.js"

lint:
	./node_modules/eslint/bin/eslint.js -c ./.eslintrs.js src/

run:
	npm run start

pushd ../frontend
pnpm build
rm -rf ../backend/web
cp -a dist ../backend/web

package main

type  Resource struct {}

func (r Resource) Get(name string) ([]byte){
   return MustAsset(name)
}
func (r Resource) GetDir(name string) ([]string, error){
   return AssetDir(name)
}

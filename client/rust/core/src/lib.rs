
#[macro_use]
extern crate serde_derive;

pub mod utils;

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
    #[test]
    fn test_stringmap() {
        use utils::{StringMap, StringMapValue};
        let mut map = StringMap::new();
        map.insert(String::from("abc"), StringMapValue::String(String::from("def")));
        assert_eq!(map["abc"], StringMapValue::String(String::from("def")));
        /*let mut str_map = StringMap::new();
        let map_box = Box::new(map);
        str_map.insert(String::from("abc"), StringMapItem::Box(map_box));
        let val = &str_map["abc"];
        match val {
            StringMapItem::Box(ans) =>  {
                assert_eq!(ans["abc"],StringMapItem::String(String::from("def")));
            },
            _ => {}
        } ;*/
    }
}

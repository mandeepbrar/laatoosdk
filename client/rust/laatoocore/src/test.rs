#[cfg(test)]
mod tests {
    use utils::{StringMap, StringMapValue};
    use redux::{Store, Action, Reducer};
    use event::{EventListener, Event};
    use platform::{Platform, SuccessCallback, ErrorCallback};
    use http;
    use application::Application;
    use std::any::Any;

    struct TestPlatform {

    }

    impl Platform for TestPlatform {
        fn http_call(&self, svc: String, method: http::HttpMethod, req: http::HttpRequest, suc: SuccessCallback, err: ErrorCallback) {

        }
    }

    #[derive(Clone, Debug)]
    struct TestData {
        testdata: Vec<i32>,
    }
    impl Store for TestData {
        fn initialize(&self) {
            //TestData{testdata: vec![]}
        }
        fn get_id(&self) -> &'static str {
            "Test Data"
        }
        fn as_any(&self) -> &dyn Any {
            self
        }

    }

/*    pub struct TestListener {

    }

    impl EventListener for TestListener {
        fn on_event(&self, evt: &Event) {
            println!("Event recieved {:?}", evt);
        }
    }
*/
    impl Reducer for TestData {
        fn reduce(&mut self, action: &Action) -> Result<bool, String> {
            match (*action).as_any().downcast_ref::<TestStoreAction>() {
                Some(act) =>  {
                    println!("Hello World{:?}", act);
                },
                None => panic!("Wrong action type!"),
            };
            /*let act = action as &TestStoreAction;
            match act {
                TestStoreAction::Add(val) => {
                    self.testdata.push(*val);
                }
            }*/
            println!("reduced");
            Ok(true)
        }

    }

    #[derive(Debug)]
    enum TestStoreAction {
        Add(i32),
    }

    impl Action for TestStoreAction {
        fn get_type(&self)->&'static str {
            return "TestStoreAction";
        }
        fn as_any(&self) -> &dyn Any {
            self
        }
    }

    #[test]
    fn store_works() {
        let mut app = create_application();
        let str = Box::new(TestData{testdata: vec![]});
        let act = TestStoreAction::Add(2);
        let str_id = str.get_id();
        app.register_store(str, act.get_type());
       // let lsr = Box::new(TestListener{});
        app.register_listener(str_id, |stor| {
            println!("event received {:?}", stor);
        });
        app.dispatch(&act);
        assert_eq!(2 + 2, 4);
    }

    fn create_application() -> Application {
        Application::new(Box::new(TestPlatform{}))
    }

    #[test]
    fn test_application() {
        let app = create_application();
        
    } 

    #[test]
    fn test_stringmap() {
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
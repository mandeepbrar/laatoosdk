use utils::StringMap;
use std::any::Any;
use std::fmt::Debug;
use error::Error;
use event::Event;


pub trait Store: Reducer  + Debug {
    fn initialize(&self);
    fn get_id(&self) -> &'static str;
    fn as_any(&self) -> &dyn Any;  
}

pub type StoreListener = fn(&Box<Store>);


/*
fn get_store<T>(evt: &Event) -> T {
    let store_evt = (*evt).as_any().downcast_ref::<StoreChangeEvent>();
    let store = store_evt.src;
    (*store).as_any().downcast_ref::<T>()
}*/


pub trait Dispatcher {
    fn dispatch(&mut self, action: &Action) -> Result<(), String>;
}

pub trait Action: Debug {
    fn get_type(&self)->&'static str;
    fn as_any(&self) -> &dyn Any;  
    //  fn get_payload(&self)->Any;
   // fn get_info(&self)->StringMap;
}


pub trait Reducer {
    /// Reduce a given state based upon an action. This won't be called externally
    /// because your application will never have a reference to the state object
    /// directly. Instead, it'll be called with you call `store.dispatch`.
    fn reduce(&mut self, &Action) -> Result<bool, String>;
}
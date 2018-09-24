use utils::StringMap;
use std::any::Any;
use std::fmt::Debug;
use error::Error;


pub trait Store: Reducer {
    fn initialize(&self);
    fn get_id(&self) -> &'static str;
}

pub trait Dispatcher {
    fn dispatch(&mut self, action: Box<Action>) -> Result<(), String>;
}

pub trait Action: Debug {
    type ActionType;
    fn get_type(&self)->&'static str;
    fn get_val(&self)-> ActionType;
    fn as_any(&self) -> &dyn Any;  
    //  fn get_payload(&self)->Any;
   // fn get_info(&self)->StringMap;
}


pub trait Reducer {
    /// Reduce a given state based upon an action. This won't be called externally
    /// because your application will never have a reference to the state object
    /// directly. Instead, it'll be called with you call `store.dispatch`.
    fn reduce(&mut self, Box<Action>) -> Result<bool, String>;
}
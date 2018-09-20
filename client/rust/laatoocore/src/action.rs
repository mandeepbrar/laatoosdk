use utils::StringMap;
use std::any::Any;

pub trait Action {
    fn get_type(&self)->&'static str;
    fn get_payload(&self)->Any;
    fn get_info(&self)->StringMap;
}

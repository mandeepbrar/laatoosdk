package utils

import (
	"log"
	"regexp"
	"testing"
)

func TestProcessTemplateRe2(t *testing.T) {
	input := `<Block className="p10 right">
          <Action   visible=$$(values.otpsent!="true")$$ id="sendotpbtn" name="sendotp" widget="button" className="submitBtn left s10">Send Otp</Action>
        </Block>
        <Field type="text" name="otp" id="otpplaceholder" visible=$$values.otpsent$$ label="javascript###replace@@@props.loginForm.otp###" placeholder='javascript###replace@@@props.loginForm.otpplaceholder###' className="text w100" module="formikforms"/> 
        <Block className="p10 right">
          <Action  visible=$$(values.otpsent=="true")$$ type="submit" widget="button" className="submitBtn left s10">Login</Action>`

	// Use TestContext from mapwriter_test.go which is in the same package 'utils'
	// It relies on embedding core.ServerContext.
	// We instantiate it simply.
	c := &TestContext{}
	// We don't populate the map because TestContext/ServerContext might be complex to init.
	// But ProcessTemplate handles empty context gracefully (returns "" or default).
	// We care about the Regex replacement which happens AFTER template execution on the buffer 'c'.
	// Wait, ProcessTemplate executes the template FIRST.
	// If the template expects specific variables, it might fail or print nothing.
	// The user input has `$$...$$` which are processed by `re2` AFTER template execution.
	// The `$$` blocks are NOT template tags (like `{{`). They are raw text to the template engine unless they overlap using delimiters.
	// The template delimiters are default `{{` and `}}`. The input has `<Block ...`.
	// So `$$` triggers represent literal text during the first pass (template execution).
	// So context values shouldn't matter for the `$$` replacement to be attempted.

	funcs := map[string]interface{}{}

	output, err := ProcessTemplate(c, []byte(input), funcs)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}

	outStr := string(output)
	log.Printf("Output: %s", outStr)

	// We expect `$$...$$` to be replaced.
	// Current regex `\$\$([.*?]*)\$\$` fails on letters.

	expectedRe := regexp.MustCompile(`javascript###replace@@@function\(values, config, target\)\{return .*?\}`)

	if !expectedRe.MatchString(outStr) {
		t.Logf("Expected regex replacement not found in output (THIS IS EXPECTED FAILURE BEFORE FIX). Got: %s", outStr)
		// We assert that it fails to confirm reproduction?
		// Or we assert it succeeds and expect the test to FAIL.
		// Standard practice: write a test that fails.
		t.Fail()
	} else {
		t.Logf("It worked? That would be unexpected with the bad regex.")
	}
}

func TestProcessTemplateComplex(t *testing.T) {
	input := `<Form title=[[settings.objectsettings[ctx.pageContext.params.item].addformtitle]] formName="additem" module="formikforms" className="p20 w100 col-xs-12 col-lg-offset-2 col-lg-8 additemform " preSubmit="preSubmitadditem" submitAction="createnewitem" >
  <Block className=" ma itemdesc">[[settings.objectsettings[ctx.pageContext.params.item].addformtext]]</Block>
  <Block className="fd fdcol ">
    <Field type="text" name="Itemname" label='Name' className="text m20 w100" module="formikforms"/>
    <Field type="bigtext" minRows="3" name="Itemdesc" label='Description' className="text m20 w100" module="formikforms"/>
    <Field type="bigtext" minRows="3" name="Itemins" label='Generation Instructions' className="text m20 w100" module="formikforms"/>
    <Field type="hidden" name="Item" value='javascript###replace@@@ctx.pageContext.params.item###' module="formikforms"/>
    <Field type="hidden" name="Application" value='javascript###replace@@@ctx.pageContext.application###' module="formikforms"/>
  </Block>
</Form>

<Form formName="login" module="formikforms" skip='true' className=" wfc ma loginform " preSubmit="preSubmitLogin" submitAction="login">
    <Block className="p20 ma loginheader">
      <Image className="block ma" src="javascript###replace@@@props.loginForm.logo###"/>
    </Block>
    <Block className="fd fdrow p20">
  
      <Block className="tb20 ps50 ma">
        <Block className="logintext tb20 f36 centered">"javascript###replace@@@props.loginForm.formtext###"</Block>
        <Field type="text" name="mobile" label="javascript###replace@@@props.loginForm.mobilelabel###" placeholder="javascript###replace@@@props.loginForm.mobileplaceholder###" className="text w100" module="formikforms"/>
        <Field type="hidden" name="otpsent" value="false" module="formikforms"/>
        <Block className="p10 right">
          <Action   visible=$$(values.otpsent!="true")$$ id="sendotpbtn" name="sendotp" widget="button" className="submitBtn left s10">Send Otp</Action>
        </Block>
        <Field type="text" name="otp" id="otpplaceholder" visible=$$values.otpsent$$ label="javascript###replace@@@props.loginForm.otp###" placeholder='javascript###replace@@@props.loginForm.otpplaceholder###' className="text w100" module="formikforms"/> 
        <Block className="p10 right">
          <Action  visible=$$(values.otpsent=="true")$$ type="submit" widget="button" className="submitBtn left s10">Login</Action>
        </Block>
      </Block>
  </Block>
</Form>`

	c := &TestContext{}
	funcs := map[string]interface{}{}
	output, err := ProcessTemplate(c, []byte(input), funcs)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}
	outStr := string(output)

	// Verify re1 replacements ([[...]])
	expectedRe1 := regexp.MustCompile(`title="javascript###replace@@@settings\.objectsettings\[ctx\.pageContext\.params\.item\]\.addformtitle###"`)
	if !expectedRe1.MatchString(outStr) {
		t.Errorf("re1 replacement failed for title. Output snippet: %s", outStr)
	}

	// Verify re2 replacements ($$...$$)
	// `visible=$$(values.otpsent!="true")$$` -> `visible="javascript###replace@@@function(values, config, target){return (values.otpsent!=&#34;true&#34;)}###"`

	expectedRe2a := regexp.MustCompile(`visible="javascript###replace@@@function\(values, config, target\)\{return \(values\.otpsent!=&#34;true&#34;\)\}###"`)
	if !expectedRe2a.MatchString(outStr) {
		t.Errorf("re2 replacement failed for visible!=true. Output snippet: %s", outStr)
	}

	expectedRe2b := regexp.MustCompile(`visible="javascript###replace@@@function\(values, config, target\)\{return values\.otpsent\}###"`)
	if !expectedRe2b.MatchString(outStr) {
		t.Errorf("re2 replacement failed for visible=values.otpsent. Output snippet: %s", outStr)
	}
}

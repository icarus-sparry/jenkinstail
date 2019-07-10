# jenkinstail
Tail a jenkins console log, exiting when job is done

The program converts CR/LF sequences to LF, rather than getting the
/consoleText file which removes all CR characters before LF.

## Examples

Get the console log for a running job

`jenkinstail http://example.com/job/daily/12345`

Add some crude timestamping. You can use `ts` from the debian `moreutils` package

`jenkinstail http://example.com/job/daily/12345 |ts`

Wait for a job to complete, then send email. Could also use banner or wall.

 `{ jenkinstail http://example.com/job/daily/12345 > /dev/null ; email -s "Finished" user@example.com } &`
 
 Look for messages to tell you when thing have happened
 
 `jenkinstail http://example.com/job/daily/12345 |grep 'Stage 1 finished'`

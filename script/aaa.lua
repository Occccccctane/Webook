---
--- Generated by EmmyLua(https://github.com/EmmyLua)
--- Created by aaa.
--- DateTime: 2024/5/10 15:33
---


EnablePrimaryMouseButtonEvents(true)
a=1.2
b=6.0
downnum = 50
down = 0
points = {
    {9630,46837,10},
    {9630,42617,10},
    {9630,38125,10},
    {9630,33390,10}
}

function OnEvent(event, arg)
    local offset
    local recovery_offset
    local downcount
    if(event== "MOUSE_BUTTON_PRESSED" and arg==1 and IsMouseButtonPressed(3)and IsKeyLockOn("capslock") ) then

        downcount=0
        if IsKeyLockOn("scrolllock") then
            down=1
        else
            down=0
        end
        repeat
            --			OutputLogMessage("%d\n",downcount)
            downcount=downcount+1
            if(downcount<=downnum and down==1) then
                MoveMouseRelative(0,b//a)
            end
            offset=math.random(6//a,14//a)
            recovery_offset = -offset
            Sleep(3)
            MoveMouseRelative(offset,offset)
            Sleep(3)
            MoveMouseRelative(recovery_offset,recovery_offset)
        until not IsMouseButtonPressed(1)
        --end
    end
end